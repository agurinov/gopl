from dataclasses import dataclass

from . import utils

def raise_exception(s: str):
	raise ValueError('oops: ' + s)

def hello_world():
	print('hello world')

def inverse_complex(s: str, b: bool, f: float, i: int) -> [str, bool, float, int]:
	return utils.ml_run(s, b, f, i)

def inverse_bools(s: str, b: bool) -> bool:
	print(f's={s} b={b}')

	if s == 'reverse':
		return not b

	return b

def inverse_float(f: float) -> float:
	print(f'f={f}')
	nf = 0 - f

	return nf

def pil(b: bytes) -> float:
	return utils.ml_image(b)

# https://docs.python.org/3/library/dataclasses.html
@dataclass
class ContractRequest:
	s: str
	b: bool
	f: float = 12.34
	i: int = 999

@dataclass
class ContractResponse:
	ml_result: float

def contract(c: ContractRequest) -> ContractResponse:
	print(f's={c.s} b={c.b} f={c.f} i={c.i}')

	if c.s == 'a':
		ml_result = 0.001
	elif c.b:
		ml_result = 0.555
	elif c.f == 100.500:
		ml_result = 100.500
	elif c.i == 100:
		ml_result = 1.123
	else:
		ml_result = -999

	return ContractResponse(ml_result=ml_result)
