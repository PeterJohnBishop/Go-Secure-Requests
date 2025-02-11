//
//  RegisterView.swift
//  automatic_fiesta
//
//  Created by Peter Bishop on 2/10/25.
//

import SwiftUI

struct RegisterView: View {
    @State var fireAuth: FireAuthViewModel = FireAuthViewModel()
    @State private var email = ""
    @State private var password = ""
    @State private var confirmPassword = ""
    @State private var qr: UIImage?
    @State private var existingUser: Bool = false

    var body: some View {
        NavigationStack{
                    VStack{
                        Spacer()
                        if let qr = qr {
                            Image(uiImage: qr)
                                .resizable()
                                .scaledToFit()
                                .frame(width: 200, height: 200)
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
                                    fireAuth.register(email: email, password: password) { result in
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
                    }
                }
    }
}

#Preview {
    RegisterView()
}
