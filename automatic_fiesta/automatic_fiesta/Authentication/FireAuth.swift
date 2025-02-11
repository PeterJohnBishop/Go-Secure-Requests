//
//  FireAuth.swift
//  automatic_fiesta
//
//  Created by Peter Bishop on 2/10/25.
//

import Foundation
import FirebaseAuth
import Observation
import PhotosUI

@Observable class FireAuthViewModel: ObservableObject {
    var email: String = ""
    var password: String = ""
    var success: Bool = false
    var response: String = ""
    var loggedIn: Bool = false
    var user: User?
    var token: String = ""
    
    private let baseURL = "http://127.0.0.1:8080"
    
    private var handle: AuthStateDidChangeListenerHandle?
    
    func register(email: String, password: String, completion: @escaping (Result<UIImage, Error>) -> Void) {
            let url = URL(string: "\(baseURL)/register")!
            var request = URLRequest(url: url)
            request.httpMethod = "POST"

            let boundary = UUID().uuidString
            request.setValue("multipart/form-data; boundary=\(boundary)", forHTTPHeaderField: "Content-Type")

            var body = Data()

            // Helper function to append form data
            func appendFormField(_ name: String, value: String) {
                body.append("--\(boundary)\r\n".data(using: .utf8)!)
                body.append("Content-Disposition: form-data; name=\"\(name)\"\r\n\r\n".data(using: .utf8)!)
                body.append("\(value)\r\n".data(using: .utf8)!)
            }

            appendFormField("email", value: email)
            appendFormField("password", value: password)

            // End boundary
            body.append("--\(boundary)--\r\n".data(using: .utf8)!)

            request.httpBody = body

            // Perform request
            URLSession.shared.dataTask(with: request) { data, response, error in
                if let error = error {
                    completion(.failure(error))
                    return
                }

                if let data = data, let image = UIImage(data: data) {
                    completion(.success(image))
                } else {
                    completion(.failure(NSError(domain: "Invalid Image", code: 0, userInfo: nil)))
                }
            }.resume()
        }
    
    func fetchFirebaseAuthToken() {
        if let user = Auth.auth().currentUser {
            user.getIDToken { token, error in
                if let error = error {
                    self.success = false
                    self.response = "Error fetching ID token: \(error.localizedDescription)"
                }
                self.success = true
                self.response = "Token retrieved!"
                self.token = token!
            }
        } else {
            self.success = false
            self.response = "No user is signed in."
        }

    }
    
    func GetCurrentUser() {
        if Auth.auth().currentUser != nil {
            self.success = true
            self.response = "Found user uid: \(String(describing: Auth.auth().currentUser?.uid))"
            self.user = Auth.auth().currentUser
        } else {
            self.success = false
            self.response = "User not found!"
        }
    }
    
    func ListenForUserState() {
        handle = Auth.auth().addStateDidChangeListener { (auth, user) in
            switch user {
            case .none:
                print("USER NOT FOUND IN CHECK AUTH STATE")
                self.loggedIn = false
            case .some(let user):
                print("FOUND: \(user.uid)!")
                self.loggedIn = true
            }
        }
    }
    
    func StopListenerForUserState() {
        if(handle != nil){
            Auth.auth().removeStateDidChangeListener(handle!)
        }
    }
    
    func SignInWithEmailAndPassword() {
        
        Auth.auth().signIn(withEmail: email, password: password) { (result, error) in
                   if error != nil {
                       self.success = false
                       self.response = error?.localizedDescription ?? ""
                   } else {
                       self.success = true
                       self.response = "Successfully signed in!"
                       self.user = Auth.auth().currentUser
                   }
               }
           
    }
    
    func SendEmailVerfication(){
        Auth.auth().currentUser?.sendEmailVerification { error in
            if error != nil {
                self.success = false
                self.response = error?.localizedDescription ?? ""
            } else {
                self.success = true
                self.response = "Email verification sent!"
            }
        }
    }

    //todo: convert existing functions to this format with completion
    func UpdateProfile(displayName: String?, photoURL: URL?, completion: @escaping (Bool, String) -> Void) {
        guard let user = Auth.auth().currentUser else {
            completion(false, "No user is currently signed in.")
            return
        }
        
        let changeRequest = user.createProfileChangeRequest()
        if let displayName = displayName {
            changeRequest.displayName = displayName
        }
        if let photoURL = photoURL {
            changeRequest.photoURL = photoURL
        }
        
        changeRequest.commitChanges { error in
            if let error = error {
                completion(false, error.localizedDescription)
            } else {
                completion(true, "Profile updated successfully!")
            }
        }
    }
    
    func UpdateEmail(newEmail: String) {
        Auth.auth().currentUser?.sendEmailVerification(beforeUpdatingEmail: newEmail) { error in
            if error != nil {
                self.success = false
                self.response = error?.localizedDescription ?? ""
            } else {
                self.success = true
                self.response = "Email updated!"
            }
        }
    }
    
    func UpdatePassword(newPassword: String) {
        Auth.auth().currentUser?.updatePassword(to: newPassword) { error in
            if error != nil {
                self.success = false
                self.response = error?.localizedDescription ?? ""
            } else {
                self.success = true
                self.response = "Password updated!"
            }
        }
    }
    
    func SendPasswordReset(){
        Auth.auth().sendPasswordReset(withEmail: email) { error in
            if error != nil {
                self.success = false
                self.response = error?.localizedDescription ?? ""
            } else {
                self.success = true
                self.response = "Password reset sent to \(self.email)!"
            }
        }
    }
    
    func DeleteCurrentUser() {
        let user = Auth.auth().currentUser

        user?.delete { error in
            if error != nil {
                self.success = false
                self.response = error?.localizedDescription ?? ""
            } else {
                self.success = true
                self.response = "User deleted!"
            }
        }
    }
    
    func SignOut(){
        let firebaseAuth = Auth.auth()
        do {
          try firebaseAuth.signOut()
            self.success = true
            self.response = "Signed out!"
        } catch let signOutError as NSError {
            self.success = false
            self.response = signOutError.description
        }
    }
}
