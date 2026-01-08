import 'package:flutter_test/flutter_test.dart';
import 'package:frontend/services/googlebooksapi_service.dart';

void main() {
group('GoogleBooksAPI', () {
  group('searchByISBN', () {
    test('returns 1984 by ISBN', () async {
      final service = GoogleBooksAPIService();
      const isbn = '0451524934';
      final result = await service.searchByISBN(isbn);

      expect(result, isNotNull);
      expect(result.title, contains('1984'));
      expect(result.authors, contains('George Orwell'));
      expect(result.isbn, contains('0451524934'));
    });
  });

  group('searchByTitle', () {
    test('returns results', () {});
  });
});
}


