import 'package:flutter/material.dart';
import 'package:frontend/services/auth_service.dart';

class LoginPage extends StatelessWidget {
  const LoginPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        children: [
          Text("Please login"),
          FilledButton.icon(
            onPressed: AuthService().signInWithGoogle,
            label: Text("Sign in with Google"),
          ),
        ],
      ),
    );
  }
}
