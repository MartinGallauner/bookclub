import 'package:flutter_test/flutter_test.dart';
import 'package:frontend/models/book.dart';
import 'package:frontend/services/library_service.dart';

void main() {
  group('LibraryService', () {
    test('New library is empty', () {
      final service = LibraryService();
      expect(service.books, hasLength(0));
    });
    test('Added books are in the collection', () {
      final service = LibraryService();
      Book book = Book(
        title: 'Nineteen Eighty-four',
        authors: ['George Orwell'],
        isbn: '0451524934',
        coverUrl: "",
      );
      service.addBook(book);

      expect(service.books, hasLength(1));
      var result = service.books[0];
      expect(result.title, contains('Nineteen Eighty-four'));
      expect(result.authors, contains('George Orwell'));
      expect(result.isbn, contains('0451524934'));
    });
  });
}
