import unittest

from internal.snowflake.snowflake import Generator, InvalidNodeID


class SnowflakeTests(unittest.TestCase):
    def test_new_invalid_node(self):
        with self.assertRaises(InvalidNodeID):
            Generator(1024)

    def test_next_id_monotonic(self):
        g = Generator(1)

        last = None
        for _ in range(10_000):
            val = g.next_id()
            if last is not None:
                self.assertGreater(val, last)
            last = val


if __name__ == "__main__":
    unittest.main()

