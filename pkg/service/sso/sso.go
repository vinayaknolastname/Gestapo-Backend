package sso

// func (auth *AuthController) SSOAuth(w http.ResponseWriter, r *http.Request) {
// 	auth.sso = &oauth2.Config{
// 		RedirectURL:  "http://localhost:8080/auth/sso-callback",
// 		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
// 		Endpoint:     google.Endpoint,
// 		ClientID:     auth.config.OAuth.WebClientId,
// 		ClientSecret: auth.config.OAuth.WebClientSecret,
// 	}
// 	url := auth.sso.AuthCodeURL(SsoOAuthString)
// 	fmt.Println(url)
// 	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
// }

// func (auth *AuthController) SSOCallback(w http.ResponseWriter, r *http.Request) {
// 	state := r.FormValue("state")
// 	code := r.FormValue("code")
// 	data, err := auth.getUserData(state, code)
// 	if err != nil {
// 		log.Fatal("error getting user data")
// 	}
// 	fmt.Fprintf(w, "Data : %s", data)
// }

// func (auth *AuthController) getUserData(state, code string) ([]byte, error) {
// 	if state != SsoOAuthString {
// 		return nil, errors.New("invalid user state")
// 	}
// 	token, err := auth.sso.Exchange(context.Background(), code)
// 	if err != nil {
// 		return nil, err
// 	}
// 	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
// 	fmt.Print(token.AccessToken)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer response.Body.Close()
// 	data, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return data, nil
// }
// func (auth *AuthController) Home(w http.ResponseWriter, r *http.Request) {
// 	render(w, "index.html")
// }

// func render(w http.ResponseWriter, t string) {
// 	partials := []string{
// 		"./static/index.html",
// 	}

// 	var templateSlice []string
// 	templateSlice = append(templateSlice, fmt.Sprintf("./static/%s", t))

// 	for i := 0; i < len(partials); i++ {
// 		templateSlice = append(templateSlice, partials...)
// 	}

// 	tmpl, err := template.ParseFiles(templateSlice...)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if err := tmpl.Execute(w, nil); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }
