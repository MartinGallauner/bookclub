import 'package:fake_cloud_firestore/fake_cloud_firestore.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:frontend/services/contact_service.dart';

void main() {
  group('ContactService', () {
    test('Add one contact to empty network.', () async {
      ContactService contactService = ContactService(FakeFirebaseFirestore());

      final contacts = await contactService.getContacts().first;

      expect(contacts, isEmpty);
    });
  });
}
