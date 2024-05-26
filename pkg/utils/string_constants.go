package utils

// Define a custom type for the context key
type contextKey string

const (
	EmailSubject string = "Sign Up OTP"
	EmailContent string = "Welcome to Gestapo !!!. Use the following OTP to complete your Sign Up procedures. OTP is valid for 5 minutes"

	Unauthorized            string     = "Unauthorized"
	AuthorizationKey        string     = "Authorization"
	AuthorizationTypeBearer string     = "bearer"
	AuthorizationPayloadKey contextKey = "authorization_payload"

	InternalServerError string = "Internal server error"
	InvalidRequest      string = "Invalid Request"
	NotFound            string = "Not Found"
	PermissionDenied    string = "Permission Denied"
	AlreadyExists       string = "Already Exists"

	PaymentCompleted string = "Payment Completed"
	PaymentPending   string = "Payment Pending"

	OrderActive    string = "Active"
	OrderCompleted string = "Completed"
	OrderCancelled string = "Cancelled"

	TrackingStatus0 int = 0
)

var TrackingTitles = []string{"Order Processed", "Order Shipped", "Order En Route", "Order Arrived"}
var TrackingSummeries = []string{"Your Order is being processed", "Your Order is Shipped", "Your Order is in Route", "Order Arrived"}
