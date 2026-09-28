package matchmaking.dominio;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

/**
 * Jogadores que entraram juntos na fila e precisam cair na mesma equipe.
 * O grupo e compartilhado pelos membros: a fila agrega os dois.
 */
public final class Grupo {
    private final String id;
    private List<Jogador> membros = new ArrayList<Jogador>();
    private boolean fechado;

    public Grupo(String id) {
        if (id == null || id.isEmpty()) {
            throw new IllegalArgumentException("grupo sem identificador");
        }
        this.id = id;
    }

    public String id() {
        return id;
    }

    public List<Jogador> membros() {
        return Collections.unmodifiableList(membros);
    }

    public int tamanho() {
        return membros.size();
    }

    void adicionar(Jogador jogador) {
        if (fechado) {
            throw new IllegalStateException("grupo ja fechado");
        }
        membros.add(jogador);
    }

    void fechar() {
        membros = List.copyOf(membros);
        fechado = true;
    }
}
