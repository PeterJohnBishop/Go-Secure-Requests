//
//  TOTPSetupView.swift
//  automatic_fiesta
//
//  Created by Peter Bishop on 2/12/25.
//

import SwiftUI
import FirebaseAuth

struct TOTPView: View {
    @State var auth: FirebaseAuth = FirebaseAuth()
    @State private var loading: Bool = false
    @State private var currentUser: User?
    @State private var code: String = ""
    @State private var validated: Bool = false
    
    var body: some View {
        NavigationStack{
            VStack{
                Text("TOTP Authentication")
                    .font(.system(size: 34, weight: .ultraLight))
                Text("Enter your code on the line below.").fontWeight(.ultraLight)
                Divider().padding()
                TextField("______", text: $code)
                    .keyboardType(.numberPad)
                    .textContentType(.oneTimeCode)
                    .onChange(of: code) { oldValue, newValue in
                        code = String(newValue.prefix(6)).filter { $0.isNumber }
                    }
                    .font(.system(size: 34, weight: .bold))
                    .multilineTextAlignment(.center)
                    .frame(height: 50)
                    .tint(.black)
                    .autocapitalization(.none)
                    .disableAutocorrection(true)
                    .padding()
                Button("Confirm", action: {
                    auth.validateTotpCode(uid: currentUser?.uid ?? "", otp: code) { result in
                        DispatchQueue.main.async {
                            switch result {
                            case .success(let valid):
                                validated = valid
                            case .failure(let error):
                                print("Error: \(error.localizedDescription)")
                            }
                        }
                    }
                })
                .navigationDestination(isPresented: $validated, destination: {
                    WelcomeView().navigationBarBackButtonHidden(true)
                })
                .fontWeight(.ultraLight)
                .foregroundColor(.black)
                .padding()
                .background(
                    RoundedRectangle(cornerRadius: 8)
                        .fill(Color.white)
                        .shadow(color: .gray.opacity(0.4), radius: 4, x: 2, y: 2)
                )
            }.onAppear{
                loading = true
                auth.GetCurrentUser(completion: {
                    user in
                    if let user = user {
                        loading = false
                        currentUser = user
                    } else {
                        print("No user is logged in.")
                    }
                })
            }
        }
    }
}

#Preview {
    TOTPView()
}
