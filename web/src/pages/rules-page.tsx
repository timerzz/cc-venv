import { ResourceFilesPage } from "./resource-files-page";

export function RulesPage() {
  return (
    <ResourceFilesPage
      descriptionKey="rules.description"
      emptyKey="rules.empty"
      guidanceIntroKey="rules.guidanceIntro"
      kind="rules"
      nameLabelKey="rules.ruleName"
      namePlaceholder="backend/api"
      titleKey="rules.title"
    />
  );
}
