import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:flutter/cupertino.dart';

import '../models/profile.dart';

class ContactService extends ChangeNotifier {
  final FirebaseFirestore _firebaseFirestore;
  List<Profile> _contacts = [];

  ContactService(this._firebaseFirestore);

  Stream<List<Profile>> getContacts() {
    return Stream.value([]);
  }

  List<Profile> get contacts => List.unmodifiable(_contacts);
}
