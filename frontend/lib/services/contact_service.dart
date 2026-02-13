import 'dart:async';

import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:flutter/cupertino.dart';
import 'package:frontend/models/connection_status.dart';

import '../models/connection.dart';

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
          uid: data['uid'],
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
    await _firebaseFirestore.collection('connections').doc().set(connection.toMap());
  }
}
