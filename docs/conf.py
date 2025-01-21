import sphinx_rtd_theme

# Configuration file for the Sphinx documentation builder.
#
# For the full list of built-in configuration values, see the documentation:
  # https://www.sphinx-doc.org/en/master/usage/configuration.html

# -- Project information -----------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#project-information

project = 'Hoverfly'
copyright = '2025, iOCO Solutions'
author = 'iOCO Solutions'

version = 'v1.10.8'
# The full version, including alpha/beta/rc tags.
release = version

# -- General configuration ---------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#general-configuration

extensions = [
    'sphinxcontrib.mermaid',
    'sphinxcontrib.jquery',
    'sphinx.ext.extlinks',
    'sphinx.ext.todo'
]

templates_path = ['_templates']

source_suffix = '.rst'

master_doc = 'index'
zip_base_url = 'https://github.com/SpectoLabs/hoverfly/releases/download/' + version + '/'

extlinks = {'zip_bundle_os_arch': (zip_base_url + 'hoverfly_bundle_%s.zip', 'zip_bundle_os_arch')}

exclude_patterns = ['_build', 'Thumbs.db', '.DS_Store']
pygments_style = 'sphinx'

# -- Options for HTML output -------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#options-for-html-output

todo_include_todos = False

html_theme = 'sphinx_rtd_theme'
html_theme_path = [sphinx_rtd_theme.get_html_theme_path()]

html_static_path = ['_static']

html_context = {
   'css_files': [
            'https://media.readthedocs.org/css/sphinx_rtd_theme.css',
            'https://media.readthedocs.org/css/readthedocs-doc-embed.css',
            '_static/theme_overrides.css',
        ],
    }


htmlhelp_basename = 'hoverflydoc'

latex_documents = [
    (master_doc, 'hoverfly.tex', u'Hoverfly Documentation',
     u'iOCO Solutions', 'manual'),
]

man_pages = [
    (master_doc, 'Hoverfly', u'Hoverfly Documentation',
     [author], 1)
]

texinfo_documents = [
    (master_doc, 'Hoverfly', u'Hoverfly Documentation',
     author, 'Hoverfly', 'API simulations for development and testing',
     'Miscellaneous'),
]
