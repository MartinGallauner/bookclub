import 'package:flutter_test/flutter_test.dart';

import 'mock_googlebooksapi_service.dart';

void main() {
  group('GoogleBooksAPI', () {
    group('searchByISBN', () {
      test('returns 1984 by ISBN', () async {
        final service = MockGooglebooksapiService();
        const isbn = '0451524934';
        final result = await service.searchByISBN(isbn);

        expect(result, isNotNull);
        expect(result.title, contains('Nineteen Eighty-four'));
        expect(result.authors, contains('George Orwell'));
        expect(result.isbn, contains('0451524934'));
        expect(result.publishedIn, equals(1949));
        expect(result.language, contains('en'));
        expect(
          result.description,
          contains(
            'Eternal warfare is the price of bleak prosperity in this satire of totalitarian barbarism.',
          ),
        );
      });
    });

    group('searchByTitle', () {
      test('returns results', () {});
    });
  });
}
