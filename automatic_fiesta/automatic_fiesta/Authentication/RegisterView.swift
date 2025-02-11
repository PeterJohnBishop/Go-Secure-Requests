//
//  RegisterView.swift
//  automatic_fiesta
//
//  Created by Peter Bishop on 2/10/25.
//

import SwiftUI

struct RegisterView: View {
    @State var fireAuth: FireAuthViewModel = FireAuthViewModel()
    @State private var email = "test@example.com"
    @State private var password = "password123"
    @State private var qr: UIImage?

    var body: some View {
        VStack {
            TextField("Email", text: $email)
                .textFieldStyle(RoundedBorderTextFieldStyle())
            SecureField("Password", text: $password)
                .textFieldStyle(RoundedBorderTextFieldStyle())

            Button("Register") {
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

            if let qr = qr {
                Image(uiImage: qr)
                    .resizable()
                    .scaledToFit()
                    .frame(width: 200, height: 200)
            }
        }
        .padding()
    }
}

#Preview {
    RegisterView()
}
