# * Currently Being Refactored for Firebase Integration *

# Go Secure Requests
A basic GO server implementing 2FA authentication, Same-Origin-Policy, and CORS 

## User Creation

![Registration](https://github.com/PeterJohnBishop/Go-Secure-Requests/blob/main/assets/1_Register.png?raw=true)

- User creates a basic account with email and password credentials, and a user account is created through Firebase Authentication.
- At the same time a secret key is generated and saved to the user Profile Document on Firestore for TOTP.

## TOTP Setup

![Setup](https://github.com/PeterJohnBishop/Go-Secure-Requests/blob/main/assets/2_TOTP_Setup.png?raw=true)

- User creation is successful when a user account is created, a user Profile doc is created, and the TOTP secret has been generated.
- A bearer token is generated from Firebase Authentication and stored in UserDefaults memory on their device for use in the Authentication header of API requests.
- The TOTP URL is converted to a QR code and displayed for the user to enable quick setup in an Authentication app.

## Login

![Login](https://github.com/PeterJohnBishop/Go-Secure-Requests/blob/main/assets/3_Login.png?raw=true)

- On login a fresh bearer token is generated from Firebase Authentication and stored in UserDefaults memory on their device for authenticating requests.

## TOTP One Time Passcode

![TOTP](https://github.com/PeterJohnBishop/Go-Secure-Requests/blob/main/assets/4_TOTP_Code.png?raw=true)

- The user has a limited amount of time to setup time based TOTP authentication in an app like Google Authenticator
- The OTP must be sent in the request as form data.
- Bearer token is verified by Firebase.
- If OTP passes verfication, a session token and a CSRF token are generated and saved to cookies with a 24h expiration. 

## Accessing a protected route
- Each request must be sent from the origin or one of the allowed cross origin domains.
- Each request must send the CSRF token as a header 'X-CSRF-Token'.
- Each request must send the session token as a cookie.
- If all of the above are satisfied and validated, a successful response can be sent.

## Logout
- On logout the values of all cookies are cleared and force expired, bearer token is revoked.