import 'dart:async';

import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:firebase_auth/firebase_auth.dart';
import 'package:frontend/models/profile.dart';

class AuthService {
  /// This class provides service functions for everything concerned with (firebase) authentication
  /// TODO By now it also holds code regarding handling of Profiles. That's something I want to extract in the future
  final FirebaseAuth firebaseAuth;
  final FirebaseFirestore firebaseFirestore;

  AuthService(this.firebaseAuth, this.firebaseFirestore);

  Future<UserCredential> signInWithGoogle() async {
    final credential = await firebaseAuth.signInWithPopup(GoogleAuthProvider());
    final user = credential.user;

    //TODO can user be null ?
    //get profile
    Profile? profile = await fetchProfile(user!.uid).first;

    if (profile == null) {
      await firebaseFirestore.collection('users').doc(user.uid).set({
        'uid': user.uid,
        'email': user.email,
        'displayName': user.displayName,
        'photoURL': user.photoURL,
        'createdAt': FieldValue.serverTimestamp(),
        'updatedAt': FieldValue.serverTimestamp(),
      });
    }
    return credential; //todo consider returning a profile instead of credentials
  }

  void signOut() {
    firebaseAuth.signOut();
  }

  Stream<User?> authStateChanges() {
    return firebaseAuth.authStateChanges();
  }

  Stream<Profile?> fetchProfile(String uid) {
    return firebaseFirestore.collection('users').doc(uid).snapshots().map((
      snapshot,
    ) {
      if (!snapshot.exists) return null;
      var data = snapshot.data()!;
      return Profile(
        uid: data['uid'],
        email: data['email'],
        displayName: data['displayName'],
        photoURL: data['photoURL'],
        updatedAt: data['updatedAt'].toDate(),
        createdAt: data['createdAt'].toDate(),
      );
    });
  }
}
