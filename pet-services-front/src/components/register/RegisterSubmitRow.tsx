import * as Form from "@radix-ui/react-form";

type RegisterSubmitRowProps = {
  isPending: boolean;
};

export default function RegisterSubmitRow({
  isPending,
}: RegisterSubmitRowProps) {
  return (
    <div className="flex justify-center">
      <Form.Submit asChild>
        <button
          type="submit"
          disabled={isPending}
          className="inline-flex h-11 w-full max-w-xs items-center justify-center rounded-full bg-linear-to-r from-teal-400 to-cyan-400 px-6 text-sm font-semibold text-white shadow-lg shadow-cyan-200 transition-opacity disabled:cursor-not-allowed disabled:opacity-70"
        >
          {isPending ? "Criando..." : "Criar conta"}
        </button>
      </Form.Submit>
    </div>
  );
}
