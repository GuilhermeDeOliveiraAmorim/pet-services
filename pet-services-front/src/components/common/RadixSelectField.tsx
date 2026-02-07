import { startTransition, useMemo, useState } from "react";
import {
  Combobox,
  ComboboxItem,
  ComboboxList,
  ComboboxProvider,
} from "@ariakit/react";
import { matchSorter } from "match-sorter";
import * as Form from "@radix-ui/react-form";
import * as Select from "@radix-ui/react-select";
import { Check, ChevronDown, Search } from "lucide-react";

export type RadixSelectOption = {
  value: string;
  label: string;
};

type RadixSelectFieldProps = {
  name: string;
  label: string;
  value: string;
  onValueChange: (value: string) => void;
  options: RadixSelectOption[];
  placeholder?: string;
  required?: boolean;
  disabled?: boolean;
  requiredMessage?: string;
  searchable?: boolean;
  searchPlaceholder?: string;
};

export default function RadixSelectField({
  name,
  label,
  value,
  onValueChange,
  options,
  placeholder = "Selecione",
  required = false,
  disabled = false,
  requiredMessage = "Obrigatório",
  searchable = false,
  searchPlaceholder = "Buscar...",
}: RadixSelectFieldProps) {
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");

  const matches = useMemo(() => {
    if (!searchable || !query.trim()) {
      return options;
    }

    const filtered = matchSorter(options, query, {
      keys: ["label", "value"],
    });

    const selectedOption = options.find((option) => option.value === value);
    if (selectedOption && !filtered.includes(selectedOption)) {
      filtered.push(selectedOption);
    }

    return filtered;
  }, [options, query, searchable, value]);

  return (
    <Form.Field className="space-y-2" name={name}>
      <div className="flex items-baseline justify-between">
        <Form.Label className="text-sm font-medium">{label}</Form.Label>
        {required ? (
          <Form.Message className="text-xs text-rose-500" match="valueMissing">
            {requiredMessage}
          </Form.Message>
        ) : null}
      </div>
      <Form.Control asChild>
        <input type="hidden" name={name} value={value} required={required} />
      </Form.Control>
      <Select.Root
        value={value}
        onValueChange={(nextValue) => {
          onValueChange(nextValue);
          setOpen(false);
        }}
        disabled={disabled}
        open={open}
        onOpenChange={(nextOpen) => {
          setOpen(nextOpen);
          if (!nextOpen) {
            setQuery("");
          }
        }}
      >
        {searchable ? (
          <ComboboxProvider
            open={open}
            setOpen={setOpen}
            resetValueOnHide
            includesBaseElement={false}
            setValue={(nextValue) => {
              startTransition(() => {
                setQuery(nextValue);
              });
            }}
          >
            <Select.Trigger
              className="flex h-11 w-full items-center justify-between gap-2 overflow-hidden rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200 disabled:cursor-not-allowed disabled:opacity-60"
              aria-label={label}
            >
              <Select.Value placeholder={placeholder} className="truncate" />
              <Select.Icon className="text-slate-400">
                <ChevronDown className="h-4 w-4" />
              </Select.Icon>
            </Select.Trigger>
            <Select.Portal>
              <Select.Content
                role="dialog"
                aria-label={label}
                position="popper"
                sideOffset={6}
                className="z-50 overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-lg"
              >
                <div className="relative border-b border-slate-100 px-3 py-2">
                  <Search className="pointer-events-none absolute left-5 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
                  <Combobox
                    autoSelect
                    placeholder={searchPlaceholder}
                    className="h-9 w-full rounded-xl border border-slate-200 bg-slate-50 pl-9 pr-3 text-sm text-slate-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
                    onBlurCapture={(event) => {
                      event.preventDefault();
                      event.stopPropagation();
                    }}
                  />
                </div>
                <ComboboxList className="max-h-60 overflow-auto p-1">
                  {matches.length === 0 ? (
                    <div className="px-3 py-2 text-sm text-slate-400">
                      Nenhum resultado
                    </div>
                  ) : null}
                  {matches.map((option) => (
                    <Select.Item
                      key={option.value}
                      value={option.value}
                      asChild
                    >
                      <ComboboxItem className="relative flex cursor-pointer select-none items-center rounded-xl px-3 py-2 text-sm text-slate-700 outline-none data-highlighted:bg-slate-100">
                        <Select.ItemText>{option.label}</Select.ItemText>
                        <Select.ItemIndicator className="absolute right-3 text-cyan-500">
                          <Check className="h-4 w-4" />
                        </Select.ItemIndicator>
                      </ComboboxItem>
                    </Select.Item>
                  ))}
                </ComboboxList>
              </Select.Content>
            </Select.Portal>
          </ComboboxProvider>
        ) : (
          <>
            <Select.Trigger
              className="flex h-11 w-full items-center justify-between gap-2 overflow-hidden rounded-2xl border border-slate-200 bg-slate-50 px-4 text-sm text-slate-700 shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200 disabled:cursor-not-allowed disabled:opacity-60"
              aria-label={label}
            >
              <Select.Value placeholder={placeholder} className="truncate" />
              <Select.Icon className="text-slate-400">
                <ChevronDown className="h-4 w-4" />
              </Select.Icon>
            </Select.Trigger>
            <Select.Portal>
              <Select.Content
                position="popper"
                sideOffset={6}
                className="z-50 overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-lg"
              >
                <Select.Viewport className="p-1">
                  {options.map((option) => (
                    <Select.Item
                      key={option.value}
                      value={option.value}
                      className="relative flex cursor-pointer select-none items-center rounded-xl px-3 py-2 text-sm text-slate-700 outline-none data-highlighted:bg-slate-100"
                    >
                      <Select.ItemText>{option.label}</Select.ItemText>
                      <Select.ItemIndicator className="absolute right-3 text-cyan-500">
                        <Check className="h-4 w-4" />
                      </Select.ItemIndicator>
                    </Select.Item>
                  ))}
                </Select.Viewport>
              </Select.Content>
            </Select.Portal>
          </>
        )}
      </Select.Root>
    </Form.Field>
  );
}
