package matchmaking.saida;

/** Os tres resultados do contrato. Quem chama so pede o texto. */
public interface Resultado {
    String texto();

    int codigoDeSaida();

    String status();

    String detalhe();

    String textoEspera();
}
