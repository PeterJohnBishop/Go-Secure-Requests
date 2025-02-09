# Go Secure Requests
A basic GO server implementing 2FA authentication, Same-Origin-Policy, and CORS 

## User Creation
- User creates a basic account with email and password credentials. 
- The password is hashed with bcrypt and basic authentication is performed by bcrypt comparison.
- At the same time a secret key is generated and saved to the user record for TOTP.
- A TOTP URL and QR code in base64 are returned in the response.

## Login
- On basic authentication a base64 token is generated and saved in a shortlived cookie.

## 2FA
- The user has a limited amount of time to setup time based TOTP authentication in an app like Google Authenticator
- The OTP must be sent in the request as form data.
- If OTP passes verfication, a session token and a CSRF token are generated and saved to cookies with a 24h expiration. 

## Accessing a protected route
- Each request must be sent from the origin or one of the allowed cross origin domains.
- Each request must send the CSRF token as a header 'X-CSRF-Token'.
- Each request must send the session token as a cookie.
- If all of the above are satisfied and validated, a successful response can be sent.

## Logout
- On logout the values of all cookies are cleared and the expiration is set to NOW.