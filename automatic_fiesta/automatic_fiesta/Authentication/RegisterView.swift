//
//  RegisterView.swift
//  automatic_fiesta
//
//  Created by Peter Bishop on 2/10/25.
//

import SwiftUI
import FirebaseAuth

struct RegisterView: View {
    @State var auth: FirebaseAuth = FirebaseAuth()
    @State private var email = ""
    @State private var password = ""
    @State private var confirmPassword = ""
    @State private var qr: UIImage?
    @State private var currentUser: User?
    @State private var userIdToken: String?
    @State private var existingUser: Bool = false
    @State private var showNext: Bool = false
    @State private var next: Bool = false

    var body: some View {
        NavigationStack{
                    VStack{
                        Spacer()
                        if let qr = qr {
                            Text("Please setup TOTP Authenticaion by scanning this QR code with an Authenticator App.").fontWeight(.ultraLight)
                                .padding()
                            Image(uiImage: qr)
                                .resizable()
                                .scaledToFit()
                                .frame(width: 200, height: 200)
                            if showNext {
                                Button("Next", action: {
                                    next.toggle()
                                }).navigationDestination(isPresented: $next, destination: {
                                    TOTPView().navigationBarBackButtonHidden(true)
                                })
                                .fontWeight(.ultraLight)
                                .foregroundColor(.black)
                                .padding()
                                .background(
                                    RoundedRectangle(cornerRadius: 8)
                                        .fill(Color.white)
                                        .shadow(color: .gray.opacity(0.4), radius: 4, x: 2, y: 2)
                                )
                            }
                            Spacer()
                        } else {
                            Text("Register").font(.system(size: 34))
                                .fontWeight(.ultraLight)
                            Divider().padding()
                            TextField("Email", text: $email)
                                .tint(.black)
                                .autocapitalization(.none)
                                .disableAutocorrection(true)
                                .padding()
                            SecureField("Password", text: $password)
                                .tint(.black)
                                .autocapitalization(.none)
                                .disableAutocorrection(true)
                                .padding()
                                .textContentType(.oneTimeCode)
                            SecureField("Confirm Password", text: $confirmPassword)
                                .tint(.black)
                                .autocapitalization(.none)
                                .disableAutocorrection(true)
                                .padding()
                                .textContentType(.oneTimeCode)
                            Button("Submit", action: {
                                if password == confirmPassword {
                                    auth.register(email: email, password: password) { result in
                                        DispatchQueue.main.async {
                                            switch result {
                                            case .success(let image):
                                                qr = image
                                            case .failure(let error):
                                                print("Error: \(error.localizedDescription)")
                                            }
                                        }
                                    }
                                }
                            })
                            .fontWeight(.ultraLight)
                            .foregroundColor(.black)
                            .padding()
                            .background(
                                RoundedRectangle(cornerRadius: 8)
                                    .fill(Color.white)
                                    .shadow(color: .gray.opacity(0.4), radius: 4, x: 2, y: 2)
                            )
                            Spacer()
                            HStack{
                                Spacer()
                                Text("I have an account.").fontWeight(.ultraLight)
                                Button("Login", action: {
                                    existingUser = true
                                }).foregroundStyle(.black)
                                    .fontWeight(.light)
                                    .navigationDestination(isPresented: $existingUser, destination: {
                                        LoginView().navigationBarBackButtonHidden(true)
                                    })
                                Spacer()
                            }
                        }
                    }.onChange(of: qr) { oldValue, newValue in
                        if newValue != nil {
                            auth.SignInWithEmailAndPassword(email: email, password: password) { result in
                                DispatchQueue.main.async {
                                    switch result {
                                    case .success(let user):
                                        currentUser = user
                                    case .failure(let error):
                                        print("Error: \(error.localizedDescription)")
                                    }
                                }
                            }
                        }
                    }
                    .onChange(of: currentUser) { oldValue, newValue in
                        if newValue != nil {
                            auth.GetIDToken(){ result in
                                DispatchQueue.main.async {
                                    switch result {
                                    case .success(let token):
                                        userIdToken = token
                                        UserDefaults.standard.set(token, forKey: "tempToken")
                                        showNext = true
                                    case .failure(let error):
                                        print("Error: \(error.localizedDescription)")
                                    }
                                }
                            }
                        }
                    }
                }
    }
}

#Preview {
    RegisterView()
}
