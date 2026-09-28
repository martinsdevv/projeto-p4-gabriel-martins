package matchmaking.dominio;

import java.util.ArrayList;
import java.util.EnumMap;
import java.util.List;
import java.util.Map;

/**
 * Cinco vagas, uma por funcao.
 * Equipe A e Equipe B sao este mesmo tipo: o que muda e o lugar na partida.
 */
public final class Equipe {
    private final EnumMap<Posicao, Vaga> vagas;

    private Equipe(EnumMap<Posicao, Vaga> vagas) {
        this.vagas = vagas;
    }

    public static Equipe montar(Map<Posicao, Jogador> ocupantes) {
        EnumMap<Posicao, Vaga> mapa = new EnumMap<Posicao, Vaga>(Posicao.class);
        for (Posicao posicao : Posicao.values()) {
            Jogador jogador = ocupantes.get(posicao);
            if (jogador == null) {
                throw new IllegalArgumentException("equipe sem " + posicao.nome());
            }
            mapa.put(posicao, new Vaga(jogador, posicao));
        }
        return new Equipe(mapa);
    }

    public int tamanho() {
        return vagas.size();
    }

    public boolean composicaoCompleta() {
        if (vagas.size() != Posicao.porEquipe()) {
            return false;
        }
        for (Posicao posicao : Posicao.values()) {
            Vaga vaga = vagas.get(posicao);
            if (vaga == null || vaga.posicao() != posicao) {
                return false;
            }
        }
        return true;
    }

    public Vaga vaga(Posicao posicao) {
        return vagas.get(posicao);
    }

    public Jogador jogadorEm(Posicao posicao) {
        return vagas.get(posicao).jogador();
    }

    public List<Vaga> vagasNaOrdem() {
        List<Vaga> lista = new ArrayList<Vaga>(Posicao.porEquipe());
        for (Posicao posicao : Posicao.values()) {
            lista.add(vagas.get(posicao));
        }
        return lista;
    }

    public List<Jogador> jogadores() {
        List<Jogador> lista = new ArrayList<Jogador>(Posicao.porEquipe());
        for (Posicao posicao : Posicao.values()) {
            lista.add(vagas.get(posicao).jogador());
        }
        return lista;
    }

    public int somaHabilidade() {
        int soma = 0;
        for (Vaga vaga : vagas.values()) {
            soma += vaga.jogador().habilidade();
        }
        return soma;
    }

    public int deslocados() {
        int total = 0;
        for (Vaga vaga : vagas.values()) {
            if (vaga.deslocada()) {
                total++;
            }
        }
        return total;
    }

    public String menorId() {
        String menor = null;
        for (Vaga vaga : vagas.values()) {
            String id = vaga.jogador().id();
            if (menor == null || id.compareTo(menor) < 0) {
                menor = id;
            }
        }
        return menor;
    }

    public String texto() {
        StringBuilder texto = new StringBuilder();
        for (Posicao posicao : Posicao.values()) {
            if (texto.length() > 0) {
                texto.append(", ");
            }
            Jogador jogador = vagas.get(posicao).jogador();
            texto.append("(").append(jogador.id()).append(", ").append(posicao.nome()).append(")");
        }
        return texto.toString();
    }
}
