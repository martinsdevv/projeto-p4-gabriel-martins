package matchmaking.avaliacao;

import matchmaking.dominio.Fila;
import matchmaking.dominio.Partida;

/** Restricao estrutural da formacao. R6 fica na politica de tolerancia. */
public interface Restricao {
    Regra regra();

    boolean satisfeita(Partida partida, Fila fila);
}
