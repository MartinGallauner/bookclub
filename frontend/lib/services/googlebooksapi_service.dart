import 'package:frontend/models/book.dart';

class GoogleBooksAPIService {
  /// This class provides service methods to call the OpenLibary.org API.

  static const String baseUrl = 'https://www.googleapis.com';

  Future<Book> searchByISBN(String isbn) async {
    return Book(title: "1984", authors: ["George Orwell"], isbn: "0451524934", coverUrl: '');
  }

}
