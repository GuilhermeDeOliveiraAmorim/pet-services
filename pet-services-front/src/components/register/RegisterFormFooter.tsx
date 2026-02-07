import Link from "next/link";

type RegisterFormFooterProps = {
  error: Error | null;
  isSuccess: boolean;
};

export default function RegisterFormFooter({
  error,
  isSuccess,
}: RegisterFormFooterProps) {
  return (
    <>
      {error ? (
        <p className="text-sm text-rose-500">
          Não foi possível criar sua conta. Verifique os dados.
        </p>
      ) : null}

      {isSuccess ? (
        <p className="text-sm text-emerald-600">
          Cadastro realizado com sucesso! Redirecionando...
        </p>
      ) : null}

      <p className="text-center text-xs text-slate-500">
        Já tem conta?{" "}
        <Link href="/login" className="text-cyan-600">
          Entrar
        </Link>
      </p>
    </>
  );
}
