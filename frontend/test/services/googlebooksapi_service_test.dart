import 'package:flutter_test/flutter_test.dart';
import 'package:frontend/services/googlebooksapi_service.dart';

void main() {
group('GoogleBooksAPI', () {
  group('searchByISBN', () {
    test('returns book data when valid', () async {
      final service = GoogleBooksAPIService();
      const isbn = '0451524934';
      final result = await service.searchByISBN(isbn);

      expect(result, "1984");



    });
    test('throws when empty', () {});
  });

  group('searchByTitle', () {
    test('returns results', () {});
  });
});
}


