from psycopg_pool import AsyncConnectionPool


def connect(database_url: str) -> AsyncConnectionPool:
    return AsyncConnectionPool(database_url, min_size=2, max_size=20, open=False)
