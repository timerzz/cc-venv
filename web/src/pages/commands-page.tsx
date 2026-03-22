import { ResourceFilesPage } from "./resource-files-page";

export function CommandsPage() {
  return (
    <ResourceFilesPage
      descriptionKey="commands.description"
      emptyKey="commands.empty"
      guidanceIntroKey="commands.guidanceIntro"
      kind="commands"
      nameLabelKey="commands.commandName"
      namePlaceholder="deploy"
      titleKey="commands.title"
    />
  );
}
