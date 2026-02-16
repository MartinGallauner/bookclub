import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:frontend/models/connection.dart';
import 'package:frontend/models/profile.dart';
import 'package:frontend/services/auth_service.dart';
import 'package:frontend/services/contact_service.dart';

class SearchService {
  final ContactService _contactService;
  final AuthService _authService;
  final FirebaseFirestore _firestore;

  SearchService(this._contactService, this._authService, this._firestore);

  Future<List<Profile>> searchByISBN(String isbn) async {
    print('search service: start search for $isbn');
    List<Profile> results = [];
    List<Connection> connections = _contactService.connections;
    print('the logged in user has number of connections ${connections.length}');
    for (var connection in connections) {
      String contactId = connection.users.firstWhere(
        (id) => id != _contactService.uid,
      );

      var doc = await _firestore
          .collection('users')
          .doc(contactId)
          .collection('library')
          .doc(isbn)
          .get();

      if (doc.exists) {
        Profile? profile = await _authService.fetchProfile(contactId);
        if (profile != null) {
          results.add(profile);
        }
      }
    }
    return Future.value(results);
  }
}
