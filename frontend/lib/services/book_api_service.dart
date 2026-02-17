import '../models/book.dart';

abstract class BookApiService {
  Future<Book> searchByISBN(String input);
}
