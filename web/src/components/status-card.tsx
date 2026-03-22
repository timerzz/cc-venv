type StatusCardProps = {
  label: string;
  title: string;
  description: string;
  tone?: "default" | "accent";
  badge?: string;
};

export function StatusCard({
  label,
  title,
  description,
  tone = "default",
  badge,
}: StatusCardProps) {
  return (
    <section
      className={
        tone === "accent"
          ? "min-h-[170px] rounded-3xl bg-[linear-gradient(180deg,#e17344_0%,#cb5d34_100%)] p-6 text-white"
          : "min-h-[170px] rounded-3xl bg-[linear-gradient(180deg,#232a30_0%,#1d2429_100%)] p-6 text-ccv-ink"
      }
    >
      <div className="text-xs uppercase tracking-[0.3em] text-white/70">{label}</div>
      <h3 className="mt-4 font-serif text-[2.1rem] font-bold leading-none">{title}</h3>
      <p className="mt-3 text-sm leading-6 text-white/75">{description}</p>
      {badge ? (
        <span className="mt-4 inline-flex min-h-8 items-center rounded-full bg-[#849d6c]/30 px-3.5 text-sm text-[#dde9cb]">
          {badge}
        </span>
      ) : null}
    </section>
  );
}
