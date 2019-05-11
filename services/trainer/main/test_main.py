import unittest

from main.main import hello_world


class HelloWorldTest(unittest.TestCase):
    def test_output(self):
        self.assertEqual(hello_world(), "Hello, world!")
