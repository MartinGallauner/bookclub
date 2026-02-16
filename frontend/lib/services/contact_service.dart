import 'dart:async';

import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:flutter/cupertino.dart';
import 'package:frontend/models/connection_status.dart';

import '../models/connection.dart';
import '../models/profile.dart';

class ContactService extends ChangeNotifier {
  final String uid;
  final FirebaseFirestore _firebaseFirestore;
  List<Connection> _connections = [];
  late StreamSubscription<QuerySnapshot> _stream;

  ContactService(this.uid, this._firebaseFirestore) {
    _stream = _firebaseFirestore
        .collection('connections')
        .where('users', arrayContains: uid)
        .snapshots()
        .listen((snapshot) {
          _connections = snapshot.docs.map((doc) {
            var data = doc.data();
            return Connection(
              uid: doc.id,
              users: List<String>.from(data['users']),
              requestedBy: data['requestedBy'],
              requestedAt: (data['requestedAt'] as Timestamp).toDate(),
              status: ConnectionStatus.values.byName(data['status']),
              rejectedAt: data['rejectedAt']?.toDate(),
              acceptedAt: data['acceptedAt']?.toDate(),

            );
          }).toList();
          notifyListeners();
        });
  }

  @override
  void dispose() {
    _stream.cancel();
    super.dispose();
  }

  List<Connection> get connections => List.unmodifiable(_connections);

  Future<void> addContact(String userToAdd) async {
    //check if target user exists

    //check if users are already connected

    //Create new connection

    Connection connection = Connection(
      requestedBy: uid,
      requestedAt: DateTime.now(),
      acceptedAt: null,
      rejectedAt: null,
      users: [uid, userToAdd],
      status: ConnectionStatus.pending,
    );
    await _firebaseFirestore
        .collection('connections')
        .doc(generateConnectionId(uid, userToAdd))
        .set(connection.toMap());
  }

  String generateConnectionId(String uid, String userToAdd) {
    List<String> sorted = [uid, userToAdd]..sort();
    return '${sorted[0]}_${sorted[1]}';
  }

  Stream<Profile?> fetchProfile(String uid) {
    return _firebaseFirestore.collection('users').doc(uid).snapshots().map((
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
