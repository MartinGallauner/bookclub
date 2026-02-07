import 'package:fake_cloud_firestore/fake_cloud_firestore.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:frontend/models/book.dart';
import 'package:frontend/services/library_service.dart';

void main() {
  group('LibraryService', () {
    test('New library is empty', () {
      final service = LibraryService('555', FakeFirebaseFirestore());
      expect(service.books, hasLength(0));
    });
    test('Added books are in the collection', () async {
      final service = LibraryService('555', FakeFirebaseFirestore());
      Book book = Book(
        title: 'Nineteen Eighty-four',
        authors: ['George Orwell'],
        isbn: '0451524934',
        coverUrl: "",
        language: 'en',
      );
      await service.addBook(book);

      expect(service.books, hasLength(1));
      var result = service.books[0];
      expect(result.title, contains('Nineteen Eighty-four'));
      expect(result.authors, contains('George Orwell'));
      expect(result.isbn, contains('0451524934'));
    });
  });
}
