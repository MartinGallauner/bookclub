import 'package:firebase_auth/firebase_auth.dart';

import '../models/profile.dart';

class AuthService {
  /// This class provides service functions for everything concerned with (firebase) authentication

  final FirebaseAuth firebaseAuth;

  AuthService(this.firebaseAuth);

  Future<UserCredential> signInWithGoogle() {
    return firebaseAuth.signInWithPopup(GoogleAuthProvider());
  }

  void signOut() {
    firebaseAuth.signOut();
  }

  Stream<User?> authStateChanges() {
    return firebaseAuth.authStateChanges();
  }

  Stream<Profile?> fetchProfile(String uid) {
    Profile profile = Profile(
      uid: 'test-uid-123',
      email: 'test@example.com',
      displayName: 'Test User',
      photoURL: 'https://example.com/photo.jpg',
    );
    return Stream<Profile>.value(profile);
  }
}
