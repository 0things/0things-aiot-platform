import { useTranslation } from 'react-i18next'
import { ContentSection } from '../components/content-section'
import { DisplayForm } from './display-form'

export function SettingsDisplay() {
  const { t } = useTranslation('settings')

  return (
    <ContentSection title={t('display.title')} desc={t('display.description')}>
      <DisplayForm />
    </ContentSection>
  )
}
