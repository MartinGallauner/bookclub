import 'package:bookclub_api/bookclub_api.dart';
import 'package:flutter/material.dart';
import 'package:frontend/main.dart';
import 'package:google_sign_in_web/web_only.dart';
import 'package:provider/provider.dart';

import '../services/auth_service.dart';

class LoginPage extends StatelessWidget {
  const LoginPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: StreamBuilder<UserProfile>(
        stream: context.read<AuthService>().authStream,
        builder: (BuildContext context, AsyncSnapshot<UserProfile> snapshot) {
          if (snapshot.hasError) {
            return Column(
              children: [
                Text("Please login"),
                renderButton(
                  configuration: GSIButtonConfiguration(
                    type: GSIButtonType.standard,
                    theme: GSIButtonTheme.filledBlue,
                    size: GSIButtonSize.large,
                  ),
                ),
                Text(
                  snapshot.error.toString(),
                  style: TextStyle(color: Colors.red),
                ),
              ],
            );
          }

          if (snapshot.hasData) {
            WidgetsBinding.instance.addPostFrameCallback((_) {
              //makes sure to run the callback only when finished loading the page.
              context.read<AppState>().login(snapshot.data!);
            });
          }

          return Column(
            children: [
              Text("Please login"),
              renderButton(
                configuration: GSIButtonConfiguration(
                  type: GSIButtonType.standard,
                  theme: GSIButtonTheme.filledBlue,
                  size: GSIButtonSize.large,
                ),
              ),
            ],
          );
        },
      ),
    );
  }
}
