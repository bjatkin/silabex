from pydantic import BaseModel
from . import Phoneme


class Connection(BaseModel):
    left: Phoneme
    right: Phoneme
    count: int


class Graph(BaseModel):
    connections: list[Connection] = []

    def add_connections(self, node_path: list[Phoneme]) -> None:
        if len(node_path) <= 1:
            return

        for conn in zip(node_path[:-1], node_path[1:]):
            self.__add_connection(conn[0], conn[1])

    def __add_connection(self, left: Phoneme, right: Phoneme):
        for conn in self.connections:
            if conn.left.phoneme == left.phoneme:
                conn.count += 1
                return

        conn = Connection(left=left, righ=right, count=1)
        self.connections.append(conn)

    def get_next(self, left: Phoneme) -> list[Connection]:
        next: list[Connection] = []
        for conn in self.connections:
            if conn.left.phoneme == left.phoneme:
                next.append(conn)

        return next

    def get_prev(self, right: Phoneme) -> list[Connection]:
        prev: list[Connection] = []
        for conn in self.connections:
            if conn.right.phoneme == right.phoneme:
                prev.append(conn)

        return prev

    def save(self, file_path: str):
        """
        Saves the graph to a file in newline-delimited json (nd_json/jsonl) file.

        Args:
            file_path (str): The path to the file where the graph will be saved.
        """

        with open(file_path, "w", encoding="utf-8") as f:
            for conn in self.connections:
                json_conn = conn.model_dump_json()
                f.write(json_conn + "\n")

    def load(self, file_path: str):
        """
        Loads the graph from an nd_json/jsonl file.

        Args:
            file_path (str): The path to the file to load graph from.

        """

        self.connections = []  # clear current connections before loading.
        with open(file_path, "r", encoding="utf-8") as f:
            for line in f:
                if line.strip():  # check if the line is not empty
                    try:
                        conn = Connection.model_validate_json(line)
                        self.connections.append(conn)
                    except Exception as e:
                        print(
                            f"Filed to create connection from json: {line.strip()}. Error: {e}"
                        )
