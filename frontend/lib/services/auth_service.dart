import 'package:firebase_auth/firebase_auth.dart';

class AuthService {
  /// This class provides service functions for everything concerned with (firebase) authentication

  Future<UserCredential> signInWithGoogle() {
    return FirebaseAuth.instance.signInWithPopup(GoogleAuthProvider());
  }

  void signOut() {
    FirebaseAuth.instance.signOut();
  }

  Stream<User?> authStateChanges() {
    return FirebaseAuth.instance.authStateChanges();
  }
}