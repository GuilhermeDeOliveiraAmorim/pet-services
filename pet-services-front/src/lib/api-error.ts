import { isAxiosError } from "axios";
import type { ProblemDetails, ProblemDetailsResponse } from "@/application";

const isProblemDetails = (value: unknown): value is ProblemDetails => {
  if (!value || typeof value !== "object") {
    return false;
  }

  const candidate = value as Partial<ProblemDetails>;

  return (
    typeof candidate.type === "string" &&
    typeof candidate.title === "string" &&
    typeof candidate.status === "number" &&
    typeof candidate.detail === "string"
  );
};

export const getApiProblemDetails = (
  error: unknown,
): ProblemDetails | undefined => {
  if (!isAxiosError(error)) {
    return undefined;
  }

  const data = error.response?.data as
    | ProblemDetailsResponse
    | ProblemDetails
    | undefined;

  if (!data) {
    return undefined;
  }

  if ("errors" in data && Array.isArray(data.errors)) {
    const first = data.errors[0];
    return isProblemDetails(first) ? first : undefined;
  }

  return isProblemDetails(data) ? data : undefined;
};

export const getApiErrorMessage = (
  error: unknown,
  fallbackMessage: string,
): string => {
  const problem = getApiProblemDetails(error);

  if (problem?.detail) {
    return problem.detail;
  }

  if (problem?.title) {
    return problem.title;
  }

  if (error instanceof Error && error.message.trim()) {
    return error.message;
  }

  return fallbackMessage;
};

export const isUnverifiedEmailError = (error: unknown): boolean => {
  const problem = getApiProblemDetails(error);

  if (!problem) {
    return false;
  }

  const title = problem.title?.toLowerCase() ?? "";
  const detail = problem.detail?.toLowerCase() ?? "";

  return (
    problem.status === 403 &&
    (title.includes("email não verificado") ||
      detail.includes("email") ||
      detail.includes("verifica"))
  );
};
