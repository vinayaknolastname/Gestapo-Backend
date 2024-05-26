package server

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/akmal4410/gestapo/pkg/grpc_api/merchant_service/db/entity"
	"github.com/akmal4410/gestapo/pkg/helpers"
	"github.com/akmal4410/gestapo/pkg/helpers/token"
	"github.com/akmal4410/gestapo/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (handler *RestServer) EditProfile(w http.ResponseWriter, r *http.Request) {
	const (
		thirtyTwoMB      = 32 << 20
		maxFileCount int = 1
	)

	// Extract the JSON data from the form
	jsonData := r.FormValue("data")
	reader := io.Reader(strings.NewReader(jsonData))

	req := new(entity.EditMerchantReq)
	err := helpers.ValidateBody(reader, req)
	if err != nil {
		handler.log.LogError("Error while ValidateBody", err)
		helpers.ErrorJson(w, http.StatusBadRequest, utils.InvalidRequest)
		return
	}

	err = r.ParseMultipartForm(thirtyTwoMB)
	if err != nil {
		handler.log.LogError("Unable to parse form", err.Error())
		helpers.ErrorJson(w, http.StatusBadRequest, utils.InvalidRequest)
		return
	}

	payload := r.Context().Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)

	files := r.MultipartForm.File["files"]
	if len(files) > maxFileCount {
		handler.log.LogError("Too many files uploaded", "Max allowed: %d", maxFileCount)
		errMsg := fmt.Sprintf("too many files uploaded. Max allowed: %s", strconv.Itoa(maxFileCount))
		helpers.ErrorJson(w, http.StatusBadRequest, errMsg)
		return
	}

	var uploadedFileKeys []string
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			handler.log.LogError("Unable to open file", err)
			helpers.ErrorJson(w, http.StatusInternalServerError, "Unable to open file")
			return
		}
		defer file.Close()

		folderPath := "profile/" + payload.UserID + "/"
		fileURL, err := handler.s3.UploadFileToS3(file, folderPath, fileHeader.Filename)
		if err != nil {
			handler.log.LogError("Error uploading file to S3", err)
			helpers.ErrorJson(w, http.StatusInternalServerError, "Error uploading file to S3")
			return
		}

		handler.log.LogInfo("File uploaded to S3 successfully", "FileURL:", fileURL)
		uploadedFileKeys = append(uploadedFileKeys, fileURL)
	}

	if len(uploadedFileKeys) != 0 {
		req.ProfileImage = uploadedFileKeys[0]
	}
	err = handler.storage.UpdateProfile(payload.UserID, req)
	if err != nil {
		handler.log.LogError("Error while UpdateProfile", err)
		helpers.ErrorJson(w, http.StatusInternalServerError, utils.InternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, "User updated successfully")
}

func (handler *RestServer) InsertProduct(w http.ResponseWriter, r *http.Request) {
	const (
		thirtyTwoMB      = 32 << 20
		maxFileCount int = 5
	)
	// Extract the JSON data from the form
	jsonData := r.FormValue("data")
	reader := io.Reader(strings.NewReader(jsonData))

	req := new(entity.AddProductReq)
	err := helpers.ValidateBody(reader, req)
	if err != nil {
		handler.log.LogError("Error while ValidateBody", err)
		helpers.ErrorJson(w, http.StatusBadRequest, utils.InvalidRequest)
		return
	}

	res, err := handler.storage.CheckDataExist("categories", "id", req.CategoryId)
	if err != nil {
		handler.log.LogError("Error while CheckCategoryExist", err)
		helpers.ErrorJson(w, http.StatusInternalServerError, utils.InternalServerError)
		return
	}
	if !res {
		err = fmt.Errorf("category doesnt exist: %s", req.CategoryId)
		handler.log.LogError("Error ", err)
		helpers.ErrorJson(w, http.StatusNotFound, err.Error())
		return
	}

	err = r.ParseMultipartForm(thirtyTwoMB)
	if err != nil {
		handler.log.LogError("Unable to parse form", err.Error())
		helpers.ErrorJson(w, http.StatusBadRequest, utils.InvalidRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		handler.log.LogError("There should be atleast one image")
		errMsg := "There should be atleast one image"
		helpers.ErrorJson(w, http.StatusBadRequest, errMsg)
		return
	}
	if len(files) > maxFileCount {
		handler.log.LogError("Too many files uploaded", "Max allowed: %d", maxFileCount)
		errMsg := fmt.Sprintf("too many files uploaded. Max allowed: %s", strconv.Itoa(maxFileCount))
		helpers.ErrorJson(w, http.StatusBadRequest, errMsg)
		return
	}

	payload, ok := r.Context().Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)
	if !ok {
		err := errors.New("unable to retrieve user payload from context")
		handler.log.LogError("Error", err)
		helpers.ErrorJson(w, http.StatusInternalServerError, utils.InternalServerError)
		return
	}
	uuId, err := uuid.NewRandom()
	if err != nil {
		handler.log.LogError("error while uuid NewRandom", err.Error())
		helpers.ErrorJson(w, http.StatusInternalServerError, utils.InternalServerError)
		return
	}

	var uploadedFileKeys []string
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			handler.log.LogError("Unable to open file", err)
			helpers.ErrorJson(w, http.StatusInternalServerError, "Unable to open file")
			return
		}
		defer file.Close()

		folderPath := filepath.Join("products", payload.UserID, uuId.String()) + "/"

		fileURL, err := handler.s3.UploadFileToS3(file, folderPath, fileHeader.Filename)
		if err != nil {
			handler.log.LogError("Error uploading file to S3", err)
			helpers.ErrorJson(w, http.StatusInternalServerError, "Error uploading file to S3")
			return
		}

		handler.log.LogInfo("File uploaded to S3 successfully", "FileURL:", fileURL)
		uploadedFileKeys = append(uploadedFileKeys, fileURL)
	}
	if len(uploadedFileKeys) != 0 {
		req.ProductImages = uploadedFileKeys
	}

	err = handler.storage.InsertProduct(payload.UserID, uuId.String(), req)
	if err != nil {
		handler.log.LogError("Error while InsertProduct", err)
		helpers.ErrorJson(w, http.StatusInternalServerError, utils.InternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, "Product added successfully")
}

func (handler *RestServer) EditProduct(w http.ResponseWriter, r *http.Request) {
	const (
		thirtyTwoMB      = 32 << 20
		maxFileCount int = 5
	)
	// Extract the JSON data from the form
	jsonData := r.FormValue("data")
	reader := io.Reader(strings.NewReader(jsonData))

	req := new(entity.EditProductReq)
	err := helpers.ValidateBody(reader, req)
	if err != nil {
		handler.log.LogError("Error while ValidateBody", err)
		helpers.ErrorJson(w, http.StatusBadRequest, err.Error())
		return
	}

	payload, ok := r.Context().Value(utils.AuthorizationPayloadKey).(*token.AccessPayload)
	if !ok {
		err := errors.New("unable to retrieve user payload from context")
		handler.log.LogError("Error", err)
		helpers.ErrorJson(w, http.StatusInternalServerError, utils.InternalServerError)
		return
	}

	id := mux.Vars(r)["id"]
	product, err := handler.storage.GetProductById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			handler.log.LogError("Error while GetProductById Not fount", err)
			helpers.ErrorJson(w, http.StatusNotFound, "Product Not found")
			return
		}
		handler.log.LogError("Error while retrieving product", err)
		helpers.ErrorJson(w, http.StatusInternalServerError, utils.InternalServerError)
		return
	}

	if *product.MerchantID != payload.UserID {
		err := errors.New("unauthorized: product does not belong to the authenticated merchant")
		handler.log.LogError("Error", err)
		helpers.ErrorJson(w, http.StatusForbidden, err.Error())
		return
	}

	if req.ClearImages {
		for _, key := range product.ProductImages {
			err := handler.s3.DeleteKey(key)
			if err != nil {
				handler.log.LogError("Error deleting file from S3", err)
				helpers.ErrorJson(w, http.StatusInternalServerError, "Error deleting file from")
				return
			}
		}
		product.ProductImages = make([]string, 0)
	}

	err = r.ParseMultipartForm(thirtyTwoMB)
	if err != nil {
		handler.log.LogError("Unable to parse form", err.Error())
		helpers.ErrorJson(w, http.StatusBadRequest, utils.InvalidRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		handler.log.LogError("There should be atleast one image")
		errMsg := "There should be atleast one image"
		helpers.ErrorJson(w, http.StatusBadRequest, errMsg)
		return
	}
	if len(files) > maxFileCount {
		handler.log.LogError("Too many files uploaded", "Max allowed: %d", maxFileCount)
		errMsg := fmt.Sprintf("too many files uploaded. Max allowed: %s", strconv.Itoa(maxFileCount))
		helpers.ErrorJson(w, http.StatusBadRequest, errMsg)
		return
	}

	var uploadedFileKeys []string
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			handler.log.LogError("Unable to open file", err)
			helpers.ErrorJson(w, http.StatusInternalServerError, "Unable to open file")
			return
		}
		defer file.Close()

		folderPath := filepath.Join("products", payload.UserID, id) + "/"

		fileURL, err := handler.s3.UploadFileToS3(file, folderPath, fileHeader.Filename)
		if err != nil {
			handler.log.LogError("Error uploading file to S3", err)
			helpers.ErrorJson(w, http.StatusInternalServerError, "Error uploading file to S3")
			return
		}

		handler.log.LogInfo("File uploaded to S3 successfully", "FileURL:", fileURL)
		uploadedFileKeys = append(uploadedFileKeys, fileURL)
	}
	uploadedFileKeys = append(uploadedFileKeys, product.ProductImages...)
	if len(uploadedFileKeys) != 0 {
		req.ProductImages = uploadedFileKeys
	}

	err = handler.storage.UpdateProduct(id, req)
	if err != nil {
		handler.log.LogError("Error while UpdateProduct", err)
		helpers.ErrorJson(w, http.StatusInternalServerError, utils.InternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, "Product updated successfully")
}
