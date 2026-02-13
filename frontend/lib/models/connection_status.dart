enum ConnectionStatus {
  pending,
  accepted,
  rejected;

  bool get isAccepted => this == ConnectionStatus.accepted;
  bool get isPending => this == ConnectionStatus.pending;


}