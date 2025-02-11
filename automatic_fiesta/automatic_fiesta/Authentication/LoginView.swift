//
//  LoginView.swift
//  automatic_fiesta
//
//  Created by Peter Bishop on 2/10/25.
//

import SwiftUI

struct LoginView: View {
    @State var fireAuth: FireAuthViewModel = FireAuthViewModel()
    @State private var email = ""
    @State private var password = ""
    @State private var newUser: Bool = false

    var body: some View {
        NavigationStack{
           VStack{
               Spacer()
               Text("Login").font(.system(size: 34))
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
               Button("Submit", action: {
                   Task{
                       fireAuth.SignInWithEmailAndPassword()
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
                   Text("I don't have an account.").fontWeight(.ultraLight)
                   Button("Register", action: {
                       newUser = true
                   }).foregroundStyle(.black)
                       .fontWeight(.light)
                       .navigationDestination(isPresented: $newUser, destination: {
                           RegisterView().navigationBarBackButtonHidden(true)
                       })
                   Spacer()
               }
           }.onAppear{
               
           }
       }
    }
}

#Preview {
    LoginView()
}
