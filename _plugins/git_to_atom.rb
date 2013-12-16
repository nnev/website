# encoding: utf-8

module Jekyll
	GIT_DIR = '/srv/git/website.git'

	class GitToRssFile < StaticFile
		def write(dest)
			super(dest) rescue ArgumentError
			true
		end
	end

	class GitToRssGenerator < Generator
		priority :low
		safe false

		def generate(site)
			require 'rss'
			require 'time'
			require 'uri'

			feed = RSS::Maker.make("atom") do |maker|
				maker.channel.link        = site.config['url']
				maker.channel.id          = tag_url(site, repo_first_change)
				maker.channel.updated     = repo_last_change
				maker.channel.title       = "#{repo_name} commits"
				maker.channel.description = "commits made to #{repo_path}"
				maker.channel.author      = site.config["author"] if site.config["author"]

				repo_recent_commits.each do |commit|
					maker.items.new_item do |item|
						item.title           = commit[:subject]
						item.description     = commit_to_description(commit)
						#item.content_encoded = commit_to_description(commit) + "\n\n\n" + commit[:changes]
						item.updated         = commit[:date]
						item.author          = commit[:author][:name]
						item.id              = tag_url(site, commit[:date], commit[:guid])
					end
				end
			end

			path = fix_path(site.config['git_to_rss_file_path'])
			dir = File.dirname(path)
			name = File.basename(path)

			FileUtils.mkdir_p(File.join(site.dest, dir))
			File.open(File.join(site.dest, path), "w") { |f| f.write(feed) }

			site.static_files << Jekyll::GitToRssFile.new(site, site.dest, dir, name)
		end

		private

		def repo_name
			File.basename(repo_path)
		end

		def repo_path
			`cd #{GIT_DIR} && git rev-parse --show-toplevel`.chomp
		end

		def repo_last_change
			Time.rfc2822(`cd #{GIT_DIR} && git log -1 --format=%aD`)
		end

		def repo_first_change
			Time.rfc2822(`cd #{GIT_DIR} && git log --format=%aD | tail -n1`)
		end

		def repo_recent_commits
			guids        = `cd #{GIT_DIR} && git log --max-count=10 --format=%H`.split("\n")
			author_names = `cd #{GIT_DIR} && git log --max-count=10 --format=%aN`.split("\n")
			author_mails = `cd #{GIT_DIR} && git log --max-count=10 --format=%aE`.split("\n")
			dates        = `cd #{GIT_DIR} && git log --max-count=10 --format=%aD`.split("\n")
			subjects     = `cd #{GIT_DIR} && git log --max-count=10 --format=%s`.split("\n")
			bodies       = `cd #{GIT_DIR} && git log --max-count=10 --format=%x00%b`.split("\0")
			files        = `cd #{GIT_DIR} && git log --max-count=10 --name-status --format=%x00`.split("\0")

			(0...guids.size).map do |i|
				{
					guid:		guids[i],
					author:	{ name: author_names[i], mail: author_mails[i] },
					date:		Time.rfc2822(dates[i]),
					subject: subjects[i],
					body:		bodies[i].strip,
					files:	 files[i].strip,
					#changes: `cd #{GIT_DIR} && git diff HEAD~#{i+1} HEAD~#{i} || echo ""`
				}
			end
		end

		def commit_to_description(commit)
			desc = [:subject, :body, :files].map do |x|
				commit[x].strip.empty? ? nil : commit[x]
			end.compact
			desc.join("\n\n\n")
		end

		# ensures the given path starts with a slash and does not end in one
		# If the given string is empty or nil, it returns /gitcommits.atom
		def fix_path(path)
			path ||= ""
			path = "/#{path}" unless path.start_with?("/")
			path = "#{path}gitcommits.atom" if path.end_with?("/")
			return path
		end

		# returns the absolute URL of the feed file, assuming the config
		# options are set correctly.
		def abs_url(site)
			url = site.config['url']
			url = url[0..-2] if site.config['url'].end_with?("/")
			url + fix_path(site.config['git_to_rss_file_path'])
		end

		# generates a tag URL used as unique identifiers for the feed and
		# its entries.
		def tag_url(site, time, extra = nil)
			uri = URI(abs_url(site))
			t = "tag:#{uri.host},#{time.strftime("%Y-%m-%d")}:#{uri.path}"
			t << "/#{extra}" if extra
			t
		end
	end
end
