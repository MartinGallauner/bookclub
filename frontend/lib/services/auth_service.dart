import 'dart:async';

import 'package:bookclub_api/bookclub_api.dart';
import 'package:frontend/models/profile.dart';
import 'package:frontend/services/api_client.dart';
import 'package:google_sign_in/google_sign_in.dart';

class AuthService {
  /// This class provides service functions for everything concerned with (firebase) authentication
  /// TODO By now it also holds code regarding handling of Profiles. That's something I want to extract in the future

  final ApiClient apiClient;
  final GoogleSignIn googleSignIn;

  AuthService(this.apiClient, this.googleSignIn);

  Stream<UserProfile> get authStream {
    return googleSignIn.authenticationEvents
        .where((event) => event is GoogleSignInAuthenticationEventSignIn)
        .map((event) => event as GoogleSignInAuthenticationEventSignIn)
        .asyncMap((event) => _handleSignedInAccount(event.user));
  }

  Future<UserProfile> _handleSignedInAccount(GoogleSignInAccount account) async {
    //extract token google-auth-token
    final idToken = account.authentication.idToken;

    //call the backend with google-auth-token
    var loginRequest = GoogleLoginRequest((b) => b.idToken = idToken);
    var loginResponse = await apiClient.authApi.loginWithGoogle(googleLoginRequest: loginRequest);

    AuthResponse authResponse = loginResponse.data!;
    //Store backend token
    apiClient.setToken(authResponse.token);

    //return UserProfile from AuthResponse
    return authResponse.user;
  }

  void signOut() {
  }


  Future<Profile?> fetchProfile(String uid) async {
    return null;


  }
}
