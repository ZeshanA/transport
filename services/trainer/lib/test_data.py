import numpy as np
import unittest

from lib.data import merge_np_tuples


class TestMergeNumpyTuples(unittest.TestCase):
    def testTwoFullLists(self):
        a = (np.array([1, 2]), np.array([3, 4]))
        b = (np.array([5, 6]), np.array([7, 8]))
        expected = (np.array([1, 2, 5, 6]), np.array([3, 4, 7, 8]))
        actual = merge_np_tuples(a, b)
        for exp, act in zip(expected, actual):
            self.assertEqual(exp.tolist(), act.tolist())
