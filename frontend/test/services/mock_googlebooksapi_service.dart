import 'package:frontend/models/book.dart';
import 'package:frontend/services/book_api_service.dart';

class MockGooglebooksapiService implements BookApiService {
  @override
  Future<Book> searchByISBN(String isbn) async {
    Book book = Book(
      title: 'Nineteen Eighty-four',
      authors: ['George Orwell'],
      isbn: '0451524934',
      publishedIn: 1949,
      coverUrl: "",
      language: "en",
      description: 'Eternal warfare is the price of bleak prosperity in this satire of totalitarian barbarism.'

    );
    return book;
  }
}
