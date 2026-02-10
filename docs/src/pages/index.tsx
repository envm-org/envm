import type {ReactNode} from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';

import styles from './index.module.css';

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
		<header className={clsx('hero hero--primary', styles.heroBanner)}>
			<div className='container'>
				<Heading as='h1' className='hero__title'>
					{siteConfig.title}
				</Heading>
				<p className='hero__subtitle'>{siteConfig.tagline}</p>
				<p className={styles.heroDescription}>
					Manage, encrypt, and sync environment variables across your team. Stop
					sharing secrets in Slack and start using ENVM.
				</p>
				<div className={styles.buttons}>
					<Link
						className='button button--secondary button--lg'
						to='/docs/intro'
					>
						Get Started →
					</Link>
					<Link
						className={clsx('button button--lg', styles.buttonOutline)}
						href='https://github.com/envm-org/envm'
					>
						View on GitHub
					</Link>
				</div>
				<div className={styles.installCommand}>
					<code>go install github.com/envm-org/envm@latest</code>
				</div>
			</div>
		</header>
	)
}

function HighlightSection() {
	return (
		<section className={styles.highlights}>
			<div className='container'>
				<div className={styles.highlightsGrid}>
					<div className={styles.highlightItem}>
						<span className={styles.highlightNumber}>256-bit</span>
						<span className={styles.highlightLabel}>AES Encryption</span>
					</div>
					<div className={styles.highlightItem}>
						<span className={styles.highlightNumber}>5+</span>
						<span className={styles.highlightLabel}>Languages Supported</span>
					</div>
					<div className={styles.highlightItem}>
						<span className={styles.highlightNumber}>100%</span>
						<span className={styles.highlightLabel}>Open Source</span>
					</div>
				</div>
			</div>
		</section>
	)
}

export default function Home(): ReactNode {
  const {siteConfig} = useDocusaurusContext();
  return (
		<Layout
			title='Secure Env Management'
			description='ENVM is a secure environment variable management and sync tool. Manage, encrypt, and share env variables across your team without exposing sensitive data.'
		>
			<HomepageHeader />
			<main>
				<HighlightSection />
				<HomepageFeatures />
			</main>
		</Layout>
	)
}
