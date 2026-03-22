import type { ReactNode } from "react";

type ModalProps = {
  open: boolean;
  title: string;
  description?: string;
  children: ReactNode;
  footer?: ReactNode;
  onClose: () => void;
};

export function Modal({ open, title, description, children, footer, onClose }: ModalProps) {
  if (!open) {
    return null;
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/55 px-4 py-6">
      <div
        aria-hidden="true"
        className="absolute inset-0"
        onClick={onClose}
      />
      <section className="relative z-10 w-full max-w-lg rounded-[1.75rem] border border-white/10 bg-[#171c20] p-6 text-ccv-ink shadow-[0_32px_96px_rgba(0,0,0,0.45)]">
        <header>
          <div className="text-xs uppercase tracking-[0.35em] text-stone-400">CCV Web</div>
          <h3 className="mt-3 font-serif text-3xl font-bold">{title}</h3>
          {description ? <p className="mt-3 text-sm leading-7 text-slate-400">{description}</p> : null}
        </header>

        <div className="mt-6">{children}</div>

        {footer ? <footer className="mt-6 flex flex-col gap-3 sm:flex-row sm:justify-end">{footer}</footer> : null}
      </section>
    </div>
  );
}
