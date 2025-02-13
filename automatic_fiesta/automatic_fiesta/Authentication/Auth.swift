//
//  Auth.swift
//  automatic_fiesta
//
//  Created by Peter Bishop on 2/12/25.
//

import Foundation
import FirebaseAuth
import Observation
import PhotosUI

@Observable class FirebaseAuth {
    
    private let baseURL = "http://127.0.0.1:8080"

    func register(email: String, password: String, completion: @escaping (Result<UIImage, Error>) -> Void) {
        
            let url = URL(string: "\(baseURL)/register")!
            var request = URLRequest(url: url)
            request.httpMethod = "POST"

            let boundary = UUID().uuidString
            request.setValue("multipart/form-data; boundary=\(boundary)", forHTTPHeaderField: "Content-Type")
            var body = Data()
            func appendFormField(_ name: String, value: String) {
                body.append("--\(boundary)\r\n".data(using: .utf8)!)
                body.append("Content-Disposition: form-data; name=\"\(name)\"\r\n\r\n".data(using: .utf8)!)
                body.append("\(value)\r\n".data(using: .utf8)!)
            }
            appendFormField("email", value: email)
            appendFormField("password", value: password)
            body.append("--\(boundary)--\r\n".data(using: .utf8)!)

            request.httpBody = body

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
    
    func SignInWithEmailAndPassword(email: String, password: String, completion: @escaping (Result<User, Error>) -> Void) {
        
        Auth.auth().signIn(withEmail: email, password: password) { (result, error) in
            if let error = error {
                completion(.failure(error))
                return
            }
            if let user = Auth.auth().currentUser {
                completion(.success(user))
            }
        }
    }
    
    func GetCurrentUser(completion: @escaping (User?) -> Void) {
        let user = Auth.auth().currentUser
        completion(user)
    }
    
    func GetIDToken(completion: @escaping (Result<String, Error>) -> Void) {
        if let user = Auth.auth().currentUser {
            user.getIDToken { token, error in
                if let error = error {
                    completion(.failure(error))
                    return
                }
                if let token = token {
                    completion(.success(token))
                }
            }
        }
    }
    
    func validateBasicAuth(uid: String, completion: @escaping (Result<Bool, Error>) -> Void) {
            let url = URL(string: "\(baseURL)/verify")!
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

            appendFormField("uid", value: uid)

            // End boundary
            body.append("--\(boundary)--\r\n".data(using: .utf8)!)

            request.httpBody = body

            // Perform request
            URLSession.shared.dataTask(with: request) { _, response, error in
                if let error = error {
                    completion(.failure(error))
                    return
                }

                if response != nil {
                    completion(.success(true))
                } 
            }.resume()
        }
    
    func validateTotpCode(uid: String, otp: String, completion: @escaping (Result<Bool, Error>) -> Void) {
        let url = URL(string: "\(baseURL)/totp")!
        var request = URLRequest(url: url)
        request.httpMethod = "POST"

        let boundary = UUID().uuidString
        request.setValue("multipart/form-data; boundary=\(boundary)", forHTTPHeaderField: "Content-Type")

        // Retrieve the JWT token from UserDefaults (or your preferred storage)
        guard let token = UserDefaults.standard.string(forKey: "tempToken") else {
            completion(.failure(NSError(domain: "AuthError", code: 401, userInfo: [NSLocalizedDescriptionKey: "Missing authentication token"])))
            return
        }

        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization") // ✅ Add Auth Header

        var body = Data()

        // Helper function to append form data
        func appendFormField(_ name: String, value: String) {
            body.append("--\(boundary)\r\n".data(using: .utf8)!)
            body.append("Content-Disposition: form-data; name=\"\(name)\"\r\n\r\n".data(using: .utf8)!)
            body.append("\(value)\r\n".data(using: .utf8)!)
        }

        appendFormField("uid", value: uid)
        appendFormField("otp", value: otp)

        // End boundary
        body.append("--\(boundary)--\r\n".data(using: .utf8)!)

        request.httpBody = body
        request.setValue("\(body.count)", forHTTPHeaderField: "Content-Length") // ✅ Add Content-Length header

        // Perform request
        URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                completion(.failure(error))
                return
            }

            guard let httpResponse = response as? HTTPURLResponse else {
                completion(.failure(NSError(domain: "InvalidResponse", code: -1, userInfo: nil)))
                return
            }

            if (200...299).contains(httpResponse.statusCode) {
                completion(.success(true)) // ✅ Check for success HTTP status codes
            } else {
                let errorMessage = "Server returned status code \(httpResponse.statusCode)"
                completion(.failure(NSError(domain: "APIError", code: httpResponse.statusCode, userInfo: [NSLocalizedDescriptionKey: errorMessage])))
            }
        }.resume()
    }

    
}
