import { ResourceFilesPage } from "./resource-files-page";

export function AgentsPage() {
  return (
    <ResourceFilesPage
      descriptionKey="agents.description"
      emptyKey="agents.empty"
      guidanceIntroKey="agents.guidanceIntro"
      kind="agents"
      nameLabelKey="agents.agentName"
      namePlaceholder="reviewer"
      titleKey="agents.title"
    />
  );
}
