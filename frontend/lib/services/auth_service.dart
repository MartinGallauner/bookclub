import 'dart:async';

import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:firebase_auth/firebase_auth.dart';
import 'package:frontend/models/profile.dart';

class AuthService {
  /// This class provides service functions for everything concerned with (firebase) authentication

  final FirebaseAuth firebaseAuth;
  final FirebaseFirestore firebaseFirestore;

  AuthService(this.firebaseAuth, this.firebaseFirestore);

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
    return firebaseFirestore.collection('users')
        .doc(uid)
        .snapshots()
        .map((snapshot) {
      if (!snapshot.exists) return null;
      var data = snapshot.data()!;
      return Profile(
        uid: data['uid'],
        email: data['email'],
        displayName: data['displayName'],
        photoURL: data['photoURL'],
      );
    });
  }
}
