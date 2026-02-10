import type {ReactNode} from 'react';
import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
	title: string
	icon: string
	description: ReactNode
}

const FeatureList: FeatureItem[] = [
	{
		title: 'Secure by Default',
		icon: '🔐',
		description: (
			<>
				End-to-end AES-256 encryption for all your secrets. Role-based access
				control, comprehensive audit logging, and secure key management keep
				your sensitive data protected.
			</>
		),
	},
	{
		title: 'Team Collaboration',
		icon: '👥',
		description: (
			<>
				Share environment variables with your team without exposing sensitive
				values. Admins can grant access to secrets that team members can use but
				never see in plain text.
			</>
		),
	},
	{
		title: 'CLI & Web GUI',
		icon: '⌨️',
		description: (
			<>
				Powerful command-line interface built with Go for maximum performance,
				plus an intuitive web dashboard for managing variables visually.
			</>
		),
	},
	{
		title: 'Multi-Environment',
		icon: '🌍',
		description: (
			<>
				Seamlessly manage dev, staging, and production environments. Support for
				environment inheritance and variable overrides makes configuration a
				breeze.
			</>
		),
	},
	{
		title: 'Version Control',
		icon: '📜',
		description: (
			<>
				Track every change to your environment variables with full history.
				Rollback to any previous version instantly when things go wrong.
			</>
		),
	},
	{
		title: 'Import & Export',
		icon: '🔄',
		description: (
			<>
				Compatible with existing tools. Import from <code>.env</code>, JSON, or
				YAML files. Export to any format for maximum flexibility and easy
				migration.
			</>
		),
	},
]

function Feature({ title, icon, description }: FeatureItem) {
	return (
		<div className={clsx('col col--4')}>
			<div className={styles.featureCard}>
				<div className={styles.featureIcon}>{icon}</div>
				<div className='padding-horiz--md'>
					<Heading as='h3' className={styles.featureTitle}>
						{title}
					</Heading>
					<p className={styles.featureDescription}>{description}</p>
				</div>
			</div>
		</div>
	)
}

export default function HomepageFeatures(): ReactNode {
  return (
		<section className={styles.features}>
			<div className='container'>
				<div className={styles.sectionHeader}>
					<Heading as='h2' className={styles.sectionTitle}>
						Everything You Need for Env Management
					</Heading>
					<p className={styles.sectionSubtitle}>
						ENVM provides a complete solution for managing environment variables
						across your entire development lifecycle.
					</p>
				</div>
				<div className='row'>
					{FeatureList.map((props, idx) => (
						<Feature key={idx} {...props} />
					))}
				</div>
			</div>
		</section>
	)
}
