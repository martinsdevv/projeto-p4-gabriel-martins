package matchmaking.saida;

public final class ErroEntrada implements Resultado {
    private final String codigo;
    private final String mensagem;

    public ErroEntrada(String codigo, String mensagem) {
        this.codigo = codigo;
        this.mensagem = mensagem;
    }

    public String codigo() {
        return codigo;
    }

    public String mensagem() {
        return mensagem;
    }

    @Override
    public String texto() {
        return "status: ERRO_ENTRADA\n"
                + "codigo: " + codigo + "\n"
                + "mensagem: " + mensagem + "\n";
    }

    @Override
    public int codigoDeSaida() {
        return 1;
    }

    @Override
    public String status() {
        return "ERRO_ENTRADA";
    }

    @Override
    public String detalhe() {
        return codigo;
    }

    @Override
    public String textoEspera() {
        return "-";
    }
}
